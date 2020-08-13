rm -rf *.lof fastdb_data iavl_data mpt_data
echo 3 > /proc/sys/vm/drop_caches
go run main.go --typ=mpt --initNumber=100000 >> test_mpt_10w.log 
go run main.go --typ=iavl --initNumber=100000 >> test_iavl_10w.log 
go run main.go --typ=fastdb --initNumber=100000 >> test_fastdb_10w.log
zip -rf *.log fastdb_data iavl_data mpt_data 10w.zip






