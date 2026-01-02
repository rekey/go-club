package mover

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rekey/go-club/common"
	"github.com/rekey/go-club/env"
)

func init() {
	DownloadDir := env.DownloadDir
	common.CreateDir(DownloadDir)
	common.CreateDir(env.DownloadResultsDir)
}

func copyAndRemove(src, dst string, srcInfo os.FileInfo) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制内容
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		// 复制失败时清理
		dstFile.Close()
		os.Remove(dst)
		return err
	}

	// 确保写入磁盘
	dstFile.Sync()

	// 复制文件属性
	os.Chtimes(dst, srcInfo.ModTime(), srcInfo.ModTime())

	// 关闭文件后再删除源文件
	srcFile.Close()
	dstFile.Close()

	// 删除源文件
	if err := os.Remove(src); err != nil {
		// 记录错误但不返回，因为复制成功了
		fmt.Printf("Warning: failed to remove source file %s: %v\n", src, err)
	}

	return nil
}

func MoveFile(src, dst string) error {
	// 检查源文件是否存在
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("source stat error: %w", err)
	}

	if !srcInfo.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	// 创建目标目录
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	err = os.Rename(src, dst)
	// 先尝试直接 rename
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	// 如果失败，使用复制+删除
	return copyAndRemove(src, dst, srcInfo)
}
