import React, { Component, Fragment } from "react";

import { withStyles } from "@material-ui/core/styles";
import AppBar from "@material-ui/core/AppBar";
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";

import StatesPage from "./pages/states";
import FixturesPage from "./pages/fixtures";
import DebugPage from "./pages/debug";

import { ApiProvider } from "./api";

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
        );
    }
}


export default withStyles(styles)(App);
