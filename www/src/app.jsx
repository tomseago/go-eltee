import React, { Component, Fragment } from "react";
import { Provider } from "react-redux";

import { withStyles } from "@material-ui/core/styles";
import AppBar from "@material-ui/core/AppBar";
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";
import Grid from "@material-ui/core/Grid";

import { ErrorBoundary } from "./common/error";

import StatesPage from "./pages/states";
import FixturesPage from "./pages/fixtures";
import DebugPage from "./pages/debug";

import { ApiProvider } from "./api";
import store from "./data";

const styles = {    
};

class App extends Component {
    state = {
        selectedTab: 0,
    };

    handleChange = (event, value) => {
        this.setState({ selectedTab: value });
    }

    debugPage() {
        return <DebugPage />;
    }

    render() {
        const { selectedTab } = this.state;

        return (
            <ErrorBoundary>
                <Provider store={store}>
                    <ApiProvider>
                        <AppBar position="static">
                            <Tabs value={selectedTab} onChange={this.handleChange}>
                                <Tab label="States" />
                                <Tab label="Fixtures" />
                                <Tab label="Debug" />
                            </Tabs>
                        </AppBar>
                        {selectedTab === 0 && <StatesPage />}
                        {selectedTab === 1 && <FixturesPage />}
                        {selectedTab === 2 && <DebugPage />}
                    </ApiProvider>
                </Provider>
            </ErrorBoundary>
        );
    }
}


export default withStyles(styles)(App);
