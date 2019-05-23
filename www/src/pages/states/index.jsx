import React, { Fragment, useState, useEffect } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";

import { withStyles } from "@material-ui/core/styles";
import Button from "@material-ui/core/Button";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
// import Divider from "@material-ui/core/Divider";

import Log from "../../lib/logger";
import Loading from "../../common/loading";
import ErrorComp, { ErrorBoundary } from "../../common/error";

// import { StateNames, ApiCall } from "../../api";
import proto from "../../api/api_pb";
import { callStateNames } from "../../data/actions";

const styles = theme => ({
    root: {
        width: "100%",
        maxWidth: 360,
        backgroundColor: theme.palette.background.paper,
    },
});

function StateListImpl(props) {
    const { names, updatedAt, isLoading, lastError, dispatch } = props;

    const [selectedName, setSelectedName] = useState(null);

    useEffect(() => {
        // Log.info(callStateNames);
        // Log.warn(cbStateNames);
        // cbStateNames();
        // //callStateNames();
        dispatch(callStateNames());
    }, []);

    if (isLoading) {
        return <Loading />;
    }

    if (lastError) {
        return <ErrorComp>{lastError}</ErrorComp>;
    }    


    function itemClick(name) {
        setSelectedName(name);
    }

    const stateItems = names.map((name) => {
        return (
            <ListItem
                button
                selected={selectedName === name}
                onClick={() => itemClick(name)}
                key={name}
            >
                <ListItemText primary={name} />
            </ListItem>
        );
    });

    return (
        <List>
            {stateItems}
        </List>
    );
}
StateListImpl.propTypes = {
    names: PropTypes.arrayOf(PropTypes.string).isRequired,
};

function mapToStateListProps(state) {
    return state.worldStates;
}

// const mapDispatchToStateListProps = (dispatch) => {
//     cbStateNames: callStateNames,
// };

const StateList = connect(
    mapToStateListProps,
)(withStyles(styles)(StateListImpl));

function StatesPageImpl(props) {
    const { dispatch } = props;

    return (
        <ErrorBoundary>
            <Fragment>
                <StateList />
                <Button variant="contained" onClick={() => dispatch(callStateNames())}>Refresh</Button>
            </Fragment>
        </ErrorBoundary>
    );
}


export default connect()(StatesPageImpl);