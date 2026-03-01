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

從 GitHub Releases 頁面下載 `RadioListUpdater.exe`。

## 執行 RadioListUpdater

執行 `RadioListUpdater.exe`。

## 從原始碼編譯與執行




## License

MIT License