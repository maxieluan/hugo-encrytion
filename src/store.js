import { configureStore, createSlice } from "@reduxjs/toolkit";

export const passwordSlice = createSlice({
    name: "password",
    initialState: {
        value: ""
    },
    reducers: {
        setPassword: (state, action) => {
            state.value = action.payload;
        }
    }
})

export const localeSlice = createSlice({
    name: "locale",
    initialState: {
        value: "zh-cn"
    },
    reducers: {
        setLocale: (state, action) => {
            state.value = action.payload;
        }
    }
})

export const setPassword = passwordSlice.actions.setPassword;
export const setLocale = localeSlice.actions.setLocale;

export default configureStore({
    reducer: {
        password: passwordSlice.reducer,
        locale: localeSlice.reducer
    }
})