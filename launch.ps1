
gulp clean

Start-Process ".\node_modules\.bin\gulp" "buildAndWatch"

Start-Process "goapp" "serve ./app-engine/"
