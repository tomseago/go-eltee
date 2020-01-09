import React, { Fragment, useState, useEffect } from "react";

import PropTypes from "prop-types";
import { connect } from "react-redux";

import { withStyles } from "@material-ui/core/styles";
import Paper from "@material-ui/core/Paper";
import Button from "@material-ui/core/Button";
import Typography from "@material-ui/core/Typography";

import proto from "../api/api_pb";
import Log from "../lib/logger";
import { doCall } from "../data/actions";


const styles = theme => ({
    root: {
        width: "100%",
        maxWidth: 360,
        backgroundColor: theme.palette.background.paper,
    },

    paper: {
        // ...theme.mixins.gutters(),
        // paddingTop: theme.spacing.unit * 2,
        // paddingBottom: theme.spacing.unit * 2,
    },
});

function PingWidgetImpl(props) {
    const { classes, dispatch, successAt, message } = props;

    function doPing() {
        dispatch(doCall("ping", new proto.StringMsg(),
            (stateCopy, resp) => {
            stateCopy.pingSuccessAt = Date.now();
            stateCopy.pingMessage = resp.getVal();
            },
            (stateCopy, error) => {
            stateCopy.pingSuccessAt = -1;
            stateCopy.pingMessage = error;
            }))
    }

    return (
        <div>
            <Button variant="contained" onClick={() => doPing()}>Ping</Button>
            <Typography>{message}</Typography>
        </div>
    );
}

export default connect(state => ({
    successAt: state.pingSuccessAt || 0,
    message: state.pingMessage,
}))(withStyles(styles)(PingWidgetImpl));
