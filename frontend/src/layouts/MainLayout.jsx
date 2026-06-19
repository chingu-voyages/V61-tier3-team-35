// layouts/MainLayout.jsx

import { Outlet } from "react-router-dom";
import Header from "../components/Header";
import Footer from "../components/Footer";

export default function MainLayout() {
    return (
        <div>
            <Header />
            <main className="min-h-[75vh]">
                <Outlet />
            </main>
            <Footer />
        </div>
    );
}