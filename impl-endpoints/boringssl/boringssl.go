package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"golang.org/x/crypto/cryptobyte"
)

var verboseMode = flag.Bool("verbose", false, "")

// extractOneECHConfig reads the base64-encoded ECHConfigList from the file at
// `echConfigListPath`. It extracts the first ECHConfig, writes it to a
// temporary file, and returns the path to that temporary file.
func extractOneECHConfig(echConfigListPath string) (string, error) {
	b64Encoded, err := ioutil.ReadFile(echConfigListPath)
	if err != nil {
		return "", err
	}
	data, err := base64.StdEncoding.DecodeString(string(b64Encoded))
	if err != nil {
		return "", err
	}
	reader := cryptobyte.String(data)
	var echConfig cryptobyte.String
	if !reader.ReadUint16LengthPrefixed(&echConfig) ||
		!reader.Empty() {
		return "", errors.New("failed to decode ECHConfigList")
	}
	tmpFile, err := ioutil.TempFile("/tmp", "*")
	if err != nil {
		return "", err
	}
	// Assume that the list contained exactly one ECHConfig.
	_, err = tmpFile.Write(echConfig)
	if err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

// extractECHServerKey reads the base64-encoded server ECH key, parses
// tls-interop-runner's ad-hoc key format, extracts just the private key, writes
// the key to a temporary file, and returns the path to that temporary file.
func extractECHServerKey(echServerKeyPath string) (string, error) {
	b64Encoded, err := ioutil.ReadFile(echServerKeyPath)
	if err != nil {
		return "", err
	}
	data, err := base64.StdEncoding.DecodeString(string(b64Encoded))
	if err != nil {
		return "", err
	}
	reader := cryptobyte.String(data)
	var key, config cryptobyte.String
	if !reader.ReadUint16LengthPrefixed(&key) ||
		!reader.ReadUint16LengthPrefixed(&config) ||
		!reader.Empty() {
		return "", errors.New("failed to decode secret key")
	}
	tmpFile, err := ioutil.TempFile("/tmp", "*")
	_, err = tmpFile.Write(key)
	if err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

func main() {
	cmd := exec.Command("/bin/sh", "/setup-routes.sh")
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	role, ok := os.LookupEnv("ROLE")
	if !ok {
		log.Fatalf("ROLE not found in env")
	}
	testCase, ok := os.LookupEnv("TESTCASE")
	if !ok {
		log.Fatalln("TESTCASE not found in env")
	}
	log.Println("Running BoringSSL server.")
	log.Println("Role:", role)
	log.Println("Test case:", testCase)

	// Convert the pregenerated ECHConfig and key into friendly formats.
	tmpECHConfigPath, err := extractOneECHConfig("/test-inputs/ech_configs")
	if err != nil {
		log.Fatalln("could not extract ECHConfig:", err)
	}
	tmpECHServerKeyPath, err := extractECHServerKey("/test-inputs/ech_key")
	if err != nil {
		log.Fatalln("could not extract ECH server key:", err)
	}

	// Start constructing the command to run the BoringSSL server.
	cmd = exec.Command("bssl", "server", "-accept", "4433")

	switch role {
	case "client":
		log.Fatalln("ROLE=client is not supported")
	case "server":
		switch testCase {
		case "dc":
			cmd.Args = append(cmd.Args, "-subcert", "/test-inputs/dc.txt",
				"-cert", "/test-inputs/example.crt",
				"-key", "/test-inputs/example.key")
		case "ech-accept":
			cmd.Args = append(cmd.Args, "-echconfig", tmpECHConfigPath,
				"-echconfig-key", tmpECHServerKeyPath,
				"-cert", "/test-inputs/example.crt",
				"-key", "/test-inputs/example.key")
		case "ech-reject":
			cmd.Args = append(cmd.Args, "-echconfig", tmpECHConfigPath,
				"-echconfig-key", tmpECHServerKeyPath,
				"-cert", "/test-inputs/client-facing.crt",
				"-key", "/test-inputs/client-facing.key")
		default:
			log.Fatalf("TESTCASE=%s not supported\n", testCase)
		}
	default:
		log.Fatalf("ROLE=%s is not supported\n", role)
	}

	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	err = cmd.Run()
	if err != nil {
		log.Fatalln("bssl server failed:", err)
	}
}
