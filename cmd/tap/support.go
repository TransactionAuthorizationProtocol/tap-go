package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	didcomm "github.com/notabene-id/go-didcomm"
	"github.com/notabene-id/go-didcomm/softkey"
)

// readInput reads message/body input: "-" or "" for stdin, "@file" for a file,
// otherwise the literal string.
func readInput(flagVal string) ([]byte, error) {
	switch {
	case flagVal == "" || flagVal == "-":
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("read stdin: %w", err)
		}
		return data, nil
	case strings.HasPrefix(flagVal, "@"):
		data, err := os.ReadFile(flagVal[1:])
		if err != nil {
			return nil, fmt.Errorf("read file: %w", err)
		}
		return data, nil
	default:
		return []byte(flagVal), nil
	}
}

// buildClient builds a DIDComm client from a keys.json (go-didcomm KeyMaterial)
// and optional comma-separated DID-document override files.
func buildClient(keyFile, didDocPaths string) (*didcomm.Client, error) {
	data, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("read key file: %w", err)
	}
	var km didcomm.KeyMaterial
	if err := json.Unmarshal(data, &km); err != nil {
		return nil, fmt.Errorf("parse key file: %w", err)
	}
	store, err := softkey.New(&km)
	if err != nil {
		return nil, err
	}

	resolver, overrides := didcomm.DefaultResolver()
	if didDocPaths != "" {
		for _, path := range strings.Split(didDocPaths, ",") {
			path = strings.TrimSpace(path)
			if path == "" {
				continue
			}
			docData, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("read DID document %s: %w", path, err)
			}
			var doc didcomm.DIDDocument
			if err := json.Unmarshal(docData, &doc); err != nil {
				return nil, fmt.Errorf("parse DID document %s: %w", path, err)
			}
			overrides.Store(&doc)
		}
	}
	return didcomm.NewClient(resolver, store), nil
}
