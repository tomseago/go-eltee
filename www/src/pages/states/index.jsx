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
import { findApiOp } from "../../api/ops";
import { maybeCallStateNames } from "../../data/actions";

import Loading from "../../common/loading";
import ErrorComp, { ErrorBoundary } from "../../common/error";
import CpList from "../control_points";

// import { StateNames, ApiCall } from "../../api";
// import proto from "../../api/api_pb";
// import { callStateNames } from "../../data/actions";

// import { ensureStateNames } from "../../api/ops";

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
    const { names, namesOp, dispatch } = props;

    useEffect(() => {
        dispatch(maybeCallStateNames(namesOp));
    }, []);

    if (namesOp.isLoading) {
        return <Loading />;
    }

    if (namesOp.lastError) {
        return <ErrorComp>{namesOp.lastError}</ErrorComp>;
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

    // return (
    //     <List>
    //         {stateItems}
    //     </List>
    // );
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
        names: state.stateNames,
        namesOp: findApiOp(state, "stateNames"),
    };
}

// const mapDispatchToStateListProps = (dispatch) => {
//     cbStateNames: callStateNames,
// };

const StateList = connect(
    mapToStateListProps,
)(withStyles(styles)(StateListImpl));

// //////////////////////////////////////////////////////////////////////

function StatesPageImpl(props) {
    const { classes, namesOp, dispatch } = props;

    return (
        <ErrorBoundary>
            <Fragment>
                <Paper className={classes.paper}>
                    <StateList />
                </Paper>
                <Button variant="contained" onClick={() => dispatch(maybeCallStateNames(namesOp))}>Refresh</Button>
            </Fragment>
        </ErrorBoundary>
    );
}


export default connect(state => ({
    namesOp: findApiOp(state, "stateNames"),
}))(withStyles(styles)(StatesPageImpl));
