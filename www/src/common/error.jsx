import React, { Component, Fragment } from "react";
import PropTypes from "prop-types";

import Log from "../lib/logger";

// import CircularProgress from "@material-ui/core/CircularProgress";

export default function ErrorComp(props) {
    const { err } = props;

    return (
        <div>
            Error: 
            {err}
        </div>
    );
}

ErrorComp.propTypes = {
    err: PropTypes.any,
};


export class ErrorBoundary extends Component {
    constructor(props) {
        super(props);
        this.state = { hasError: false };
    }

    static getDerivedStateFromError(err) {
        return { hasError: true };
    }

    componentDidCatch(err, info) {
        Log.error("===== Error Boundary Hit ====");
        Log.error(err);
    }

    render() {
        const { hasError } = this.state;
        const { children } = this.props;

        if (hasError) {
            return (
                <Fragment>
                    <h1>Oh Noes!</h1>
                    <h2>Some sort of really bad error has happened. Oops.</h2>
                </Fragment>
            );
        }

        return children;
    }
}

ErrorBoundary.propTypes = {
    children: PropTypes.element.isRequired,
};
