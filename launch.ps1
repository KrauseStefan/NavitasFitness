
.\node_modules\.bin\gulp clean

Start-Process ".\node_modules\.bin\gulp" "buildAndWatch"

Start-Process "goapp" "serve ./app-engine/"

cd ipn-simulator
Start-Process "npm start"
cd ..
