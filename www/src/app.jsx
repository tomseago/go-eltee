import React, { Component, Fragment } from "react";
import { Provider } from "react-redux";

import { MuiThemeProvider, createMuiTheme, withStyles } from "@material-ui/core/styles";
import AppBar from "@material-ui/core/AppBar";
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";
import Grid from "@material-ui/core/Grid";

import styles, { theme } from "./styles";
import { ErrorBoundary } from "./common/error";


import StatesPage from "./pages/states";
import FixturesPage from "./pages/fixtures";
import DebugPage from "./pages/debug";

import { ApiProvider } from "./api";
import store from "./data";



class AppUiImpl extends Component {
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
        const { classes } = this.props;        
        const { selectedTab } = this.state;

        return (
            <ErrorBoundary>                
                <Provider store={store}>
                    <ApiProvider>
                        <div className={classes.root}>
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
                        </div>
                    </ApiProvider>
                </Provider>
            </ErrorBoundary>
        );
    }
}


const AppUi = withStyles(styles)(AppUiImpl);

export default function App(props) {
    return (
        <MuiThemeProvider theme={theme}>
            <AppUi />
        </MuiThemeProvider>
    );
}
