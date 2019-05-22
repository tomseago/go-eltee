const gulp = require("gulp");
const gulpParcel = require("gulp-parcel");

const OUT = "./root";

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




module.exports = {
    parcelProd: tParcelProd,
    parcel: tParcel,
}