import type {Metadata} from "next";
import {Geist} from "next/font/google";
import "./globals.css";
import Header from "@/app/components/header";
import LocalFont from "next/font/local";
import Starter from "@/app/components/starter";

const geistSans = Geist({
    variable: "--font-geist-sans",
    subsets: ["latin"],
});

const soyuzGrotesk = LocalFont({
    variable: "--font-soyuz-grotesk",
    src: [
        {
            path: "./../public/font/SoyuzGroteskBold.otf",
            style: "normal",
        },
    ],
})

const manrope = LocalFont({
    variable: "--font-manrope",
    src: [
        {
            path: "./../public/font/manrope-regular.otf",
            style: "normal",
            weight: "500",
        },
    ],
})

export const metadata: Metadata = {
    title: "Начало",
    description: "Краудфандинговая платформа - cбор денег для бизнеса, технологических, творческих и социальных проектов.",
};

export default function RootLayout({
                                       children,
                                   }: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="ru">
        <body className={geistSans.variable + ' ' + soyuzGrotesk.variable + ' ' + manrope.variable}
              style={{minHeight: '100vh'}}>
        <Header/>
        {children}
        <Starter/>
        </body>
        </html>
    );
}
