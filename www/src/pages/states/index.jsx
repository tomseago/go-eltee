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

import Log from "../../lib/logger";
import { findApiCall } from "../../api/ops";
import { maybeCallStateNames } from "../../data/actions";

import Loading from "../../common/loading";
import ErrorComp, { ErrorBoundary } from "../../common/error";
import CpList from "../control_points";

import PingWidget from "../../common/ping_widget";


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

function StateListImpl(props) {
    const { names, namesCall, dispatch } = props;

    useEffect(() => {
        dispatch(maybeCallStateNames());
    }, []);

    if (namesCall.isLoading) {
        return <Loading />;
    }

    if (namesCall.lastError) {
        return <ErrorComp>{namesCall.lastError}</ErrorComp>;
    }

    const stateItems = names.map(name => (
        <ExpansionPanel key={name}>
            <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                <Typography variant="h4">{name}</Typography>
            </ExpansionPanelSummary>
            <ExpansionPanelDetails>
                <CpList wsName={name} />
            </ExpansionPanelDetails>
        </ExpansionPanel>
    ));

    return (
        <Fragment>
            {stateItems}
        </Fragment>
    );
}
StateListImpl.propTypes = {
    names: PropTypes.arrayOf(PropTypes.string).isRequired,
};

function mapToStateListProps(state) {
    return {
        names: state.stateNames || [],
        namesCall: findApiCall(state, "StateNames"),
    };
}

const StateList = connect(state => ({
    names: state.stateNames || [],
    namesCall: findApiCall(state, "StateNames"),
}))(withStyles(styles)(StateListImpl));

// //////////////////////////////////////////////////////////////////////

function StatesPageImpl(props) {
    const { classes, dispatch } = props;

    return (
        <ErrorBoundary>
            <Fragment>
                <Paper className={classes.paper}>
                    <StateList />
                </Paper>
                <Button variant="contained" onClick={() => dispatch(maybeCallStateNames())}>Refresh</Button>
                <PingWidget/>
            </Fragment>
        </ErrorBoundary>
    );
}


export default connect(state => ({
}))(withStyles(styles)(StatesPageImpl));
