import "./App.css";

import { MantineProvider, Global } from "@mantine/core";
import { NotificationsProvider } from "@mantine/notifications";

import { BrowserRouter, Routes, Route } from "react-router-dom";

import Dashboard from "./routes/Dashboard";
import Login from "./routes/Login";
import Register from "./routes/Register";
import DashboardRedirect from "./components/dashboard/Redirect";
import { useEffect, useState } from "react";
import Landing from "./routes/Landing";

function App() {
    const isDark = window.localStorage.getItem("dark") === "true";

    useEffect(() => {
        document.documentElement.setAttribute(
            "data-theme",
            isDark ? "dark" : "light"
        );
    }, [isDark]);

    return (
        <>
            <Global
                styles={() => ({
                    body: {
                        color: "var(--text-color)",
                    },
                })}
            />
            <MantineProvider
                theme={{
                    colorScheme: isDark ? "dark" : "light",
                }}
            >
                <NotificationsProvider>
                    <BrowserRouter>
                        <Routes>
                            <Route path="/" element={<Landing />} />
                            <Route
                                path="/dashboard"
                                element={<DashboardRedirect />}
                            />
                            <Route
                                path="/dashboard/:handle"
                                element={<Dashboard />}
                            />
                            <Route path="/login" element={<Login />} />
                            <Route path="/register" element={<Register />} />
                        </Routes>
                    </BrowserRouter>
                </NotificationsProvider>
            </MantineProvider>
        </>
    );
}

export default App;
