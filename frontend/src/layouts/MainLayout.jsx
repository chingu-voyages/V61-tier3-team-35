import { Outlet } from "react-router-dom";
import Footer from "../components/Footer";

export default function MainLayout() {
    return (
        <div className="bg-primary h-screen overflow-hidden">
            <main className="">
                <Outlet />
            </main>
            <Footer />
        </div>
    );
}