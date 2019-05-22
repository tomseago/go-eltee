const { exec } = require("child_process");

const gulp = require("gulp");
const gulpParcel = require("gulp-parcel");

const OUT = "./root";

///////

function execCmd(cmd, cb) {
    //const cmd = `docker tag ssi:latest ettpdemo.azurecr.io/tp/ssidata:latest`;
    console.log("Running ", cmd);
    const x = exec(cmd, (err, stdout, stderr) => {
        cb(err);
    });
    x.stdout.on("data", (data) => {
        console.log(data);
    });
    x.stderr.on("data", (data) => {
        console.error(data);
    });
}

//////

function tParcelProd() {
    process.env.NODE_ENV = "production";
    return gulp.src("./src/index.html", {read:false})
        .pipe(gulpParcel({
            outDir: `${OUT}`,
            cache: false,
        }));
}

function tParcel() {
    return gulp.src("./src/index.html", {read:false})
        .pipe(gulpParcel({
            outDir: `${OUT}`,
            cache: true,
            watch: true,
        }));
}

function tProtoc(cb) {
    execCmd(`protoc -I../api ../api/api.proto --js_out=import_style=commonjs:src --grpc-web_out=import_style=commonjs,mode=grpcwebtext:src`, cb)
}


module.exports = {
    parcelProd: tParcelProd,
    parcel: tParcel,

    protoc: tProtoc,
}