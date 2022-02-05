# Toggl Time Entry Manipulator
## 概要
Toggl trackで時間を記録するためのalfred workflow。  

Togglで時間を記録するのに加えて、見積もりを記録することができる。  
見積もりはtogglで記録できないので、別途Firestoreを利用している。

https://user-images.githubusercontent.com/11168945/151706717-6aea3f00-f713-42e4-987e-919ba7ac5043.mp4

## 設定
### configファイル
configは以下のファイルを編集する。

`$HOME/Library/Application Support/Alfred/Workflow Data/com.hytssk.toggl_time_entry_manipulator/config.json`

```json 
{
    "TogglConfig": {
        "APIKey":  "xxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    },
	"FirestoreConfig": {
        "CollectionName": "private"
    },
    "WorkflowConfig": {
        "ProjectAutocompleteItems": ["[PBI]", "[Event]"],
        "ProjectAliases": [
            {"ID": 1234567, "Alias": "tag alias"}
        ],
        "TagAliases": [
            {"ID": 1111111, "Alias": "project alias"}
        ],
        "RecordEstimate": false
    }
}

```

- `TogglConfig.APIKey`: togglのapiキー
- `FirestoreConfig.CollectionName`: 見積もりを記録するためのFirestoreのコレクション名
- `WorkflowConfig.ProjectAutocompleteItems`: プロジェクト一覧にオートコンプリート用のitemを追加 (プロジェクトに共通のprefixがついているときに使う想定)
- `WorkflowConfig.ProjectAutocompleteItems`: プロジェクトのIDに対してエイリアスを設定する (例: 日本語のプロジェクト名に対して英語のエイリアスを設定)
- `WorkflowConfig.TagAutocompleteItems`: タグのIDに対してエイリアスを設定する (例: 日本語のタグ名に対して英語のエイリアスを設定)
- `WorkflowConfig.RecordEstimate`: 見積もり記録の有効・無効を切り替える
### Firestoreの認証
サービスアカウントのjsonキーを使う。  

参考: https://www.ipentec.com/document/software-google-cloud-platform-get-service-account-key

#### 見積もりを記録するために必要な設定
- jsonキーは `$HOME/Library/Application Support/Alfred/Workflow Data/com.hytssk.toggl_time_entry_manipulator/secret.json` に保存する
- configファイルのRecordEstimateをtrueにする
