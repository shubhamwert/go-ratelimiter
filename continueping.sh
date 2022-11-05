for VARIABLE in {1..100}
do
    echo VARIABLE
    echo $(curl "http://localhost:8080/get?api-key=123")
done