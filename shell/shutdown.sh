PROCESS='goweb'
PID=$(ps -e|grep goweb |grep -v grep|awk '{printf $1}')

if [ $? -eq 0 ]; then
    echo "pid:$PID"
else
    echo "进程 $PROCESS 不存在"
    exit
fi

kill -9 ${PID}

if [ $? -eq 0 ];then
    echo "kill $PROCESS success"
else
    echo "kill $PROCESS fail"
fi