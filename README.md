## JSON2EXCEL

```go
    package main

    import (
        "log"

        "github.com/user0608/json2excel"
    )
    func main(){
        var data=`[{
				"nombre":"Kevin",
				"apellido":"Saucedo",
				"edad":25,
				"estado":true    
			}]`
        req, err := json2excel.NewRequest([]byte(data))
        if err != nil {
            log.Fatalln(err)
        }
        converter := json2excel.NewJSON2ExcelConverter()
        if err := converter.SaveExcel(req, "file.xlsx"); err != nil {
            log.Fatalln(err)
        }
    }    
```