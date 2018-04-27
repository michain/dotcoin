### Hashs (散列)
通常情况下，dotcoin使用SHA-256散列，RIPEMD-160会用于生成较短的散列(例如生成比特币地址的时候)。

对字符串"hello"进行double-SHA-256散列计算的例子:

```
hello
2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824 (第一轮 sha-256)
9595c9df90075148eb06860365df33584b75bff782a510c6cd4883a419833d50 (第二轮 sha-256)
```

生成账户地址时(RIPEMD-160)会得到：

```
hello
2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824 (第一轮 使用 sha-256)
b6a9c8c230722b7c748331a8b450f05566dc7d0f (第二轮 使用 ripemd-160)
```

### Merkle Trees (Merkle树)
Merkle树是散列的二叉树。在dotcoin中，Merkle树使用SHA-256算法，是这样生成的：

```
sha256(a) sha256(b) sha256(c)
sha256(sha256(a)+sha256(b)) sha256(sha256(c)+sha256(nil))
sha256(sha256(sha256(a)+sha256(b))+sha256(sha256(c)+sha256(c)))
```
每轮都将上一轮的结果两两相接后计算，若最后剩余单个元素则复制后计算。

### Signatures (签名)
dotcoin使用椭圆曲线数字签名算法(ECDSA)对交易进行签名

公钥(in scripts) 以 04 <x> <y>的形式给出，x和y是表示曲线上点的坐标的32字节字符串。签名使用DER 编码 将 r 和 s 写入一个字节流中(因为这是OpenSSL的默认输出).