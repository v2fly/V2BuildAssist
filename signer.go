package V2BuildAssist

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/v2fly/VSign/insmgr"
	"github.com/v2fly/VSign/sign"
	"github.com/v2fly/VSign/signerVerify"
)

func RequestForSign(githubToken,
	signerOwner, signerRepo,
	project,
	projectOwner, projectRepo,
	version,
	keyPassword, keyEncrypted string) (int64, string, []byte, error) {
	data, id, err := GetReleaseFile(githubToken, projectOwner, projectRepo, version, "Release.unsigned")
	fmt.Fprintln(os.Stderr, "Getting Release")
	if err != nil {
		return 0, "", nil, err
	}
	//Sort First, this will also read all ins, if some are invalid, it will crash here
	sorted := bytes.NewBuffer(nil)
	insmgr.SortAll(bytes.NewReader(data), sorted)
	//Check version and project
	insall := insmgr.ReadAllIns(bytes.NewReader(sorted.Bytes()))
	if !signerVerify.CheckVersionAndProject(insall, version, project) {
		fmt.Fprintln(os.Stderr, "Cannot Check Constraint")
		return 0, "", nil, io.EOF
	}
	r := base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(keyEncrypted)))
	keyorig, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot Read key")
		return 0, "", nil, err
	}
	sr, err := sign.Sign(keyorig, keyPassword, sorted.Bytes())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot Sign")
		return 0, "", nil, err
	}
	sfpath := fmt.Sprintf("%v/%v.Release", project, version)
	releasebuf := bytes.NewBuffer(nil)
	releasebuf.WriteString("untrusted comment: Signed Release File\n")
	releasebuf.WriteString(string(sr))
	releasebuf.WriteString("\n")
	releasebuf.Write(sorted.Bytes())
	url, err := CreateFileIfNotExist(githubToken, signerOwner, signerRepo, sfpath, "Signed Release", releasebuf.Bytes())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot Upload Data")
		return 0, "", nil, err
	}
	return id, url, releasebuf.Bytes(), nil
}
