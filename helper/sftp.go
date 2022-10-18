package helper

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path"
)

type SftpCli struct {
	SshClient     *ssh.Client
	Client     *sftp.Client
}

func (s *SftpCli) Connect() error {
	var sftpClient *sftp.Client
	var err error
	if sftpClient, err = sftp.NewClient(s.SshClient); err != nil {
		return err
	}
	s.Client = sftpClient
	return nil
}

func (s *SftpCli) UploadFile(localFilePath string, remotePath string) error{
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("os.Open error : ", localFilePath)
		return err
	}
	defer srcFile.Close()
	var remoteDir = path.Dir(remotePath)
	_,err = s.Client.Stat(remoteDir)
	if os.IsNotExist(err){
		err = s.Client.MkdirAll(remoteDir)
		fmt.Println(err)
		return err
	}
	dstFile, err := s.Client.Create(remotePath)
	if err != nil {
		fmt.Println("sftpClient.Create error : ", localFilePath, remotePath)
		return err
	}
	defer dstFile.Close()
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localFilePath)
		return err
	}
	dstFile.Write(ff)
	fmt.Println("文件上传成功")
	return nil
}


func (s *SftpCli) DownloadFile(localFilePath string, remotePath string) error{
	srcFile, _ := s.Client.Open(remotePath) //远程
	dstFile, _ := os.Create(localFilePath) //本地
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()

	if _, err := srcFile.WriteTo(dstFile); err != nil {
		return err
	}
	fmt.Println("文件下载完毕")
	return nil
}
