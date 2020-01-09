import React, { Fragment, useState, useEffect } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";

import { withStyles } from "@material-ui/core/styles";
import Paper from "@material-ui/core/Paper";
import Button from "@material-ui/core/Button";

import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
// import Divider from "@material-ui/core/Divider";

import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import Typography from '@material-ui/core/Typography';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';

import Log from "../lib/logger";
import { findApiCall } from "../api/ops";
import { maybeCallControlPoints } from "../data/actions";

import Loading from "../common/loading";
import ErrorComp, { ErrorBoundary } from "../common/error";


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

/////////////
const fakeWidget = kind => (props) => {
    const { cp, onChange } = props;
    const [val, setVal] = useState({});

    return (
        <Typography variant="body2">
            {kind}
            {JSON.stringify(cp)}
        </Typography>
    );
};

const CpColorWidget = fakeWidget("Color");
const CpXyzWidget = fakeWidget("Xyz");
const CpEnmWidget = fakeWidget("Enm");
const CpIntensityWidget = fakeWidget("Intensity");

/////////////
function ControlPointListImpl(props) {
    const { wsName, cpList, cpOp, dispatch } = props;

    useEffect(() => {
        dispatch(maybeCallControlPoints(wsName));
    }, []);

    if (cpOp.isLoading) {
        return <Loading />;
    }

    if (cpOp.lastError) {
        return <ErrorComp>{cpOp.lastError}</ErrorComp>;
    }

    function colorChanged(cpName, val) {
        Log.info("Color CP ", cpName, val);
    }

    function xyzChanged(cpName, val) {
        Log.info("XYZ CP ", cpName, val);
    }

    function enmChanged(cpName, val) {
        Log.info("Enum CP ", cpName, val);
    }

    function intensityChanged(cpName, val) {
        Log.info("Intensity CP ", cpName, val);
    }

    Log.info("cpList ", cpList);
    const cpItems = !cpList ? null : cpList.getCpsList().map((cp) => {
        const cpName = cp.getName();

        return (
            <ListItem key={cpName}>
                <ListItemText>{cpName}</ListItemText>
                {cp.hasColor() && <CpColorWidget cp={cp.getColor().toObject()} onChange={val => colorChanged(cpName, val)} />}
                {cp.hasXyz() && <CpXyzWidget cp={cp.getXyz().toObject()} onChange={val => xyzChanged(cpName, val)} />}
                {cp.hasEnm() && <CpEnmWidget cp={cp.getEnm().toObject()} onChange={val => enmChanged(cpName, val)} />}
                {cp.hasIntensity() && <CpIntensityWidget cp={cp.getIntensity().toObject()} onChange={val => intensityChanged(cpName, val)} />}
            </ListItem>
        );
    });

    return (
        <List>
            {cpItems}
        </List>
    );
}
ControlPointListImpl.propTypes = {
//    names: PropTypes.arrayOf(PropTypes.string).isRequired,
};


export default connect((state, ownProps) => {
    const { wsName } = ownProps;
    Log.info(`State ${wsName} state.cpsByState =`, state.cpsByState)

    return {
        cpList: state.cpsByState[wsName],
        cpOp: findApiCall(state, `controlPoints:${wsName}`)
    };
})(withStyles(styles)(ControlPointListImpl));
