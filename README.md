ノベルゲームエンジン。シナリオテキストを入れて、ノベルゲームを作成できる。

## 動作例

- https://kijimad.github.io/nova/ アプリケーション
- https://github.com/kijimaD/nova/tree/main/_example ソースコード

## 記法

吉里吉里の記法を参考にした。

- `[p]`: クリック待ちにし、クリック時に表示内容をリセットする
- `[l]`: クリック待ちにし、クリック時に改行する
- `[r]`: 改行する
- `[image source="test.png"]`: 背景を表示する
- `[jump target="label1"]`: TARGETのラベルに移動する
- `[wait time="1000"]`: TIMEミリ秒操作待ちにする
- `*this_is_label`: ラベル定義。`start`ラベルを最初に読み込む
