import { MuiThemeProvider, createMuiTheme, withStyles } from "@material-ui/core/styles";

export const theme = createMuiTheme({
    palette: {
        primary: {
            main: "#1c1c1c",
        },
        secondary: {
            main: "#af8700",
        },
        type: "dark",
        background: {
            paper: "#585858",
            default: "#002b36",
        },
    },
    typography: {
        useNextVariants: true,
    },    
});

export default function styles() {
    return th => ({
        root: {
            backgroundColor: th.palette.background.default,
        },

        paper: {
            ...theme.mixins.gutters(),
            paddingTop: theme.spacing.unit * 2,
            paddingBottom: theme.spacing.unit * 2,
        },
    });
}
