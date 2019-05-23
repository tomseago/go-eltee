import React, { Fragment } from "react";
import Fade from "@material-ui/core/Fade";
import CircularProgress from "@material-ui/core/CircularProgress";

export default function Loading() {
    return (
        <div style={{ textAlign: "center" }}>
            <Fade
                in
                unmountOnExit
                style={{
                    transitionDelay: "800ms",
                }}
            >
                <CircularProgress />
            </Fade>
        </div>
    );
}

export function TextLoading() {
    return (
        <Fragment>
            Loading...
        </Fragment>
    );
}
