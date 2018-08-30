let gobr = require("gonobridge")

let goemitter = gobr.Init("go.wasm", [])

goemitter.on("handler", function (handler) {
    const express = require('express')
    const app = express()

    app.get('/', (req, res) => {
        s = 0;
        for (let i = 0; i < 1000000; i++) {
            let x = i;
            while(x != 0) {
                x &= x - 1;
                s++;
            }
        }
        res.send('Hello World from js!')
    })
    app.get('/go*', handler)

    app.listen(3000, () => console.log('Example app listening on port 3000!'))
})