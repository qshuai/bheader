### bindex
====

bindex aims to get blockheader's information and saved to the specified file.  At the same time, bindex will count the number of blockheader and output it. The blockheader will be sorted by block height. It is noted that bindex will only collect data set at the current block height and the bitcoin client should be shut down.

#### Usage: 

1. clone this repository:

   ```
   git clone https://github.com/qshuai/bindex.git $GOPATH/src/bindex
   ```

2. install the tool:

   ```
   cd $GOPATH/src/bindex
   go install
   ```

3. go ahead now:

   ```
   bindex --dbdir=~/.bitcoin/blocks/index --file=headers-mainnet.dat
   // now you will get a sorted blockheader file
   ```

