# RadioListUpdater

這是 [CRadio](https://github.com/riddleling/cradio) 的電台列表更新程式，主要是會抓取 `台北愛樂電台` 與 `港都電台` 的串流 URL，並加到 `list.json` 裡，如果 `list.json` 裡已有這兩個電台，則會更新串流 URL。

註：因為`台北愛樂電台` 與 `港都電台` 的串流 URL 包含隨機生成的「Token」以及「有效時間」，所以使用此程式更新串流 URL。


## 系統需求

- Go 1.20+ (如要自行編譯原始碼)
- yt-dlp (必要)

### 安裝 yt-dlp

- Windows：
    ```
    winget install yt-dlp.yt-dlp
    ```
- macOS：
    ```
    brew install yt-dlp
    ```
- Ubuntu/Debian：
    ```
    pip install yt-dlp
    ```

## 下載執行檔 (Windows)

從 [GitHub Releases](https://github.com/riddleling/RadioListUpdater/releases) 頁面下載 `RadioListUpdater.exe`。

## 執行 RadioListUpdater

把執行檔 `RadioListUpdater.exe` 放到 CRadio 目錄下 (裡頭有 `list.json`)，然後執行 `RadioListUpdater.exe` 就會更新電台列表 。

## 從原始碼編譯

```
git clone https://github.com/riddleling/RadioListUpdater.git
cd RadioListUpdater
go mod tidy
go build
```

## 建立 .bat 檔案

我執行的流程是先執行 `RadioListUpdater.exe` 更新電台列表，然後再執行 `cradio.exe` 來播放電台，所以可以寫一個 `run_cradio.bat` 來執行此流程：

run_cradio.bat 內容：

```
@echo off
cd /d "%~dp0"

RadioListUpdater.exe || exit /b

timeout /t 1 /nobreak >nul

cradio.exe
```

以後要用 `CRadio` 播放電台，都直接執行這個 .bat 檔案即可。


## License

MIT License
