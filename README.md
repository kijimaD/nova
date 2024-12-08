ノベルゲームエンジン。シナリオテキストを入れて、ノベルゲームを作成できる。

## 動作例

- https://kijimad.github.io/nova/ サンプルアプリ
  - https://github.com/kijimaD/nova/tree/main/_example コード
- https://kijimad.github.io/na2me/ 青空文庫ビューワ
  - https://github.com/kijimaD/na2me コード

## シナリオファイル例

```
*start
[image source="black.png"]
『吾輩は猫である』夏目漱石
[p]
[jump target="ch1"]

*ch1
一
[p]
[image source="sky.jpg"]
[wait time="500"]

吾輩は猫である。名前はまだ無い。
[p]
どこで生れたかとんと見当がつかぬ。何でも薄暗いじめじめした所でニャーニャー泣いていた事だけは記憶している。
[p]
[jump target="start"]
```

## 記法

吉里吉里の記法を参考にした。

- `[p]`: クリック待ちにし、クリック時に表示内容をリセットする
- `[l]`: クリック待ちにし、クリック時に改行する
- `[r]`: 改行する
- `[image source="test.png"]`: 背景を表示する
- `[jump target="label1"]`: TARGETのラベルに移動する
- `[wait time="1000"]`: TIMEミリ秒操作待ちにする
- `*this_is_label`: ラベル定義。`start`ラベルを最初に読み込む
