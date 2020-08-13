rm -rf *.log *data
echo 3 > /proc/sys/vm/drop_caches
go run main.go --typ=mpt    --initNumber=1000000 >> test_mpt_100w.log
go run main.go --typ=iavl   --initNumber=1000000 >> test_iavl_100w.log
go run main.go --typ=fastdb --initNumber=1000000 >> test_fastdb_100w.log
go run main.go --typ=mpts   --initNumber=1000000 >> test_mpts_100w.log
zip -r 100w.zip *.log *_data




rm -rf  *_data
echo 3 > /proc/sys/vm/drop_caches
go run main.go --typ=mpt    --initNumber=10000000 >> test_mpt_1000w.log
go run main.go --typ=iavl   --initNumber=10000000 >> test_iavl_1000w.log
go run main.go --typ=fastdb --initNumber=10000000 >> test_fastdb_1000w.log
go run main.go --typ=mpts   --initNumber=10000000 >> test_mpts_1000w.log
zip -r 1000w.zip *.log *_data




rm -rf  *_data
echo 3 > /proc/sys/vm/drop_caches
go run main.go --typ=mpt    --initNumber=100000000 >> test_mpt_10000w.log
go run main.go --typ=iavl   --initNumber=100000000 >> test_iavl_10000w.log
go run main.go --typ=fastdb --initNumber=100000000 >> test_fastdb_10000w.log
go run main.go --typ=mpts   --initNumber=100000000 >> test_mpts_10000w.log
zip -r 10000w.zip *.log *_data
