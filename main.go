package main

import (
	"encoding/json"
	"faceAPI/s3Downloader"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 顔の情報
type FaceInfo struct {
	Image string
	Match int
}

// 顔の画像目と一致率を一時的に保存する
var faceInfo = FaceInfo{
	Image: "",
	Match: -1,
}

// POSTリクエストに対する処理
func postJsonHandler(rw http.ResponseWriter, req *http.Request) {
	// リクエストの設定
	rw.Header().Set("Content-Type", "application/json")

	// メソッドの確認
	if req.Method != "POST" {
		fmt.Fprint(rw, "Method Not POST.")
		return
	}

	// データの取得
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprint(rw, err.Error())
		fmt.Println(err.Error())
		return
	}

	// データの代入
	err = json.Unmarshal(body, &faceInfo)
	if err != nil {
		fmt.Fprint(rw, err.Error())
		fmt.Println(err.Error())
		return
	}

	// 受け取ったデータを表示
	fmt.Printf("%#v\n", faceInfo)

	// クライアントへのレスポンス
	fmt.Fprint(rw, "success post data!")
}

// GETリクエストに対する処理
func getJsonHandler(rw http.ResponseWriter, req *http.Request) {
	// リクエストの設定
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// JSONエンコード
	faceJson, err := json.Marshal(faceInfo)
	if err != nil {
		fmt.Fprint(rw, err.Error())
		return
	}

	// クライアントへ取得したデータを投げる
	fmt.Fprint(rw, string(faceJson))

	// 一致率が0未満の場合は、データが来ていないのでS3から画像を取得しない
	if faceInfo.Match < 0 {
		return
	}

	// S3から顔画像を取得してくる
	imageDownload(faceInfo.Image)

	// 一時的に保存していたデータの初期化
	faceInfo.Image = ""
	faceInfo.Match = -1
}

// S3から指定した画像を取得して保存する
func imageDownload(fileName string) {
	// S3の設定
	s3Info := s3Downloader.S3DownloadInfo{
		AccessKeyId:     "アクセスキー",
		SecretAccessKey: "シークレットキー",
		Region:          "ap-northeast-1", // 東京リージョン
		BucketName:      "バケット名",
	}
	// 保存するディレクトリパス
	filePath := "./../face-app/html/images/"

	// ファイルのダウンロード
	s3Downloader.FileDownloadFromS3(s3Info, filePath, fileName)
}

func main() {
	// ハンドラの設定
	http.HandleFunc("/post", postJsonHandler)
	http.HandleFunc("/get", getJsonHandler)
	http.ListenAndServe(":9999", nil)
}
