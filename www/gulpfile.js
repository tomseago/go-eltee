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


module.exports = {
    parcelProd: tParcelProd,
}