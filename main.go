package main

import (
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/pavel-v-chernykh/keystore-go/v4"
)

type nonRand struct{}

func (r nonRand) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 1
	}

	return len(p), nil
}

func zeroing(buf []byte) {
	for i := range buf {
		buf[i] = 0
	}
}

func readPrivateKey(filepath string) []byte {
	pkPEM, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	b, _ := pem.Decode(pkPEM)
	if b == nil {
		log.Fatal("should have at least one pem block")
	}

	if b.Type != "PRIVATE KEY" {
		log.Fatal("should be a private key")
	}

	return b.Bytes
}

func readCertificate(filepath string) []byte {
	pkPEM, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	b, _ := pem.Decode(pkPEM)
	if b == nil {
		log.Fatal("should have at least one pem block")
	}

	if b.Type != "CERTIFICATE" {
		log.Fatal("should be a certificate")
	}

	return b.Bytes
}

func writeKeyStore(ks keystore.KeyStore, filename string, password []byte) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	err = ks.Store(f, password)
	if err != nil {
		panic(err)
	}
}

func main() {

	ct := time.Now()

	a := app.New()
	w := a.NewWindow("Keystore Browser")
	w.Resize(fyne.NewSize(600, 600))
	label := widget.NewLabel("open a pem file")
	pk := keystore.PrivateKeyEntry{
		CreationTime: ct,
	}
	ks := keystore.New(
		keystore.WithOrderedAliases(),
		keystore.WithCustomRandomNumberGenerator(nonRand{}),
	)
	w.SetContent(container.NewVBox(
		label,
		widget.NewButton("select certificate pem file", func() {
			dialog.ShowFileOpen(func(certUri fyne.URIReadCloser, err error) {
				log.Println("open dialog to select file: %w", certUri.URI().Path())
				certificate := readCertificate(certUri.URI().Path())
				pk.CertificateChain = []keystore.Certificate{
					{
						Type:    "X509",
						Content: certificate,
					},
				}
			}, w)
		}),
		widget.NewButton("select private key pem file", func() {
			dialog.ShowFileOpen(func(keyUri fyne.URIReadCloser, err error) {
				log.Println("open dialog to select file: %w", keyUri.URI().Path())
				privateKey := readPrivateKey(keyUri.URI().Path())
				pk.PrivateKey = privateKey
			}, w)
		}),
		widget.NewButton("generate key store", func() {
			dialog.ShowFileSave(func(uc fyne.URIWriteCloser, err error) {
				fileExt := uc.URI().Extension()
				filepath := uc.URI().Path()

				var popUp *widget.PopUp
				entry := widget.NewPasswordEntry()
				form := &widget.Form{
					Items: []*widget.FormItem{
						{Text: "Password", Widget: entry}},
					OnCancel: func() {
						log.Println("cancelled")
						os.Remove(filepath)
						popUp.Hide()
					},
					OnSubmit: func() {
						log.Println("password:", entry.Text)
						password := []byte(entry.Text)
						defer zeroing(password)
						if err := ks.SetPrivateKeyEntry("pk", pk, password); err != nil {
							panic(err)
						}

						if strings.EqualFold(".jks", fileExt) || strings.EqualFold(".jceks", fileExt) {
							writeKeyStore(ks, filepath, password)
							popUp.Hide()
						} else {
							log.Fatal("file not saved successfully due to incorrect file extension")
						}
					},
				}

				popUp = widget.NewModalPopUp(form, w.Canvas())
				popUp.Show()
			}, w)

		}),
	))
	w.ShowAndRun()
}
