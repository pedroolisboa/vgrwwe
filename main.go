package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

// Caminhos comuns para localizar o TestWe.exe
var possibleDirectories = []string{
	fmt.Sprintf("%s\\TestWe\\TestWe.exe", os.Getenv("APPDATA")),
	fmt.Sprintf("%s\\Programs\\TestWe\\TestWe.exe", os.Getenv("LOCALAPPDATA")),
	fmt.Sprintf("%s\\TestWe\\TestWe.exe", os.Getenv("PROGRAMFILES")),
	fmt.Sprintf("%s\\TestWe\\TestWe.exe", os.Getenv("PROGRAMFILES(X86)")),
}

func findTestWeExecutable() (string, error) {
	for _, dir := range possibleDirectories {
		if _, err := os.Stat(dir); err == nil {
			return dir, nil
		}
	}
	return "", fmt.Errorf("TestWe.exe não foi encontrado em nenhum dos diretórios conhecidos")
}

func modifyExecutableData(data string) string {
	// Alterar strings que indicam que é uma máquina virtual
	vmIndicators := []string{"QEMU", "VirtualBox", "VMware", "Parallels", "Hyper-V", "Xen"}

	// Substituir as ocorrências de indicadores de VM por algo "inocente"
	for _, indicator := range vmIndicators {
		re := regexp.MustCompile("(?i)" + indicator) // ignorar maiúsculas e minúsculas
		data = re.ReplaceAllString(data, "InnocentPC") // substitui as ocorrências
	}

	return data
}

func main() {
	// Encontrar o TestWe.exe
	directory, err := findTestWeExecutable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}

	// Ler o arquivo TestWe.exe
	exe, err := ioutil.ReadFile(directory)
	if err != nil {
		panic(err)
	}

	data := string(exe)

	// Modificar os dados do executável
	modifiedData := modifyExecutableData(data)

	// Renomear o arquivo original e salvar a versão modificada
	backupFileName := directory + ".bak"
	err = os.Rename(directory, backupFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao renomear o arquivo original: %v\n", err)
		os.Exit(1)
	}

	file, err := os.Create(directory)
	if err != nil {
		panic(err)
	}

	_, err = file.Write([]byte(modifiedData))
	if err != nil {
		panic(err)
	}

	_ = file.Close()

	fmt.Println("Modificações aplicadas com sucesso. O arquivo original foi salvo como", backupFileName)
}
