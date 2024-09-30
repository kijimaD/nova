ノベルゲームエンジン。シナリオテキストを入れて、ノベルゲームを作成できる。

## 動作例

- https://kijimad.github.io/nova/ アプリケーション
- https://github.com/kijimaD/nova/tree/main/_example ソースコード

## 記法

吉里吉里の記法を参考にした。

- `[p]`: 改ページクリック待ち。クリック時にページをフラッシュする
- `[l]`: クリック待ち
- `[image source="test.png"]`: 背景表示
- `[jump target="label1"]`: ラベル移動
- `[wait time="1000"]`: 操作待ち。単位はミリ秒
- `*{LABEL}`: ラベル定義。デフォルトでは`start`ラベルを最初に読み込む
