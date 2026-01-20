import type { Metadata } from "next";
// Импортируем Montserrat
import { Montserrat } from "next/font/google";
import "./globals.css";
import Header from "@/app/components/header";
import Starter from "@/app/components/starter";

// Настраиваем Montserrat (ExtraBold для заголовков)
const montserrat = Montserrat({
    subsets: ["latin", "cyrillic"],
    weight: ["400", "700", "800"], // 800 для ExtraBold
    variable: "--font-montserrat",
});

export const metadata: Metadata = {
    title: "СИПиС", // Обновим название
    description: "Система инвестиционных проектов и сборов",
};

export default function RootLayout({
                                       children,
                                   }: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="ru">
        {/* Применяем переменные шрифтов. Для основного текста используем системные шрифты (Calibri есть в sans-serif на Windows) */}
        <body className={`${montserrat.variable} font-sans`}
              style={{minHeight: '100vh', fontFamily: 'Calibri, Arial, sans-serif'}}>
        <Header/>
        {children}
        <Starter/>
        </body>
        </html>
    );
}