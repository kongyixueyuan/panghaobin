

### UTXO模型:
 >  比特币交易模型为UTXO(unspend transaction output),只记录未花费的的比特币输出，每次交易花费的是之前的一个或者多个输出
 
 > UTXO实现：
 
![image](https://github.com/drpang/readme_add_pic/blob/master/Images/UTXO.png)

比特币的交易中实际没有账户余额信息，只有输入(input)与输出(output)两部分信息。
 
 说的更直白一点，比特币的交易类似于现实世界中直接用支票进行交易。如A写给B一张100元的支票，其中A是输入，B是输出。当B需要给C付50元时，B不是去银行兑现支票然后付款给C，而是直接写两张新支票，一张50元给C，一张50元给自己，再将原来100元的支票作废，即完成一次非基于账户模型的支付。