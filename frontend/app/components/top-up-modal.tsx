'use client'
import {useState} from "react";
import {Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger} from "@/app/components/ui/dialog";
import {Button} from "@/app/components/ui/button";
import {Input} from "@/app/components/ui/input";
import {Label} from "@/app/components/ui/label";
import MessageComponent from "@/app/components/message";
import {Message} from "@/api/api";
import Spinner from "@/app/components/spinner";
import {InitPayment} from "@/api/payment";
import {PlusCircle, CreditCard} from "lucide-react";

export default function TopUpModal() {
    const [amount, setAmount] = useState<string>("1000");
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState<Message | null>(null);
    const [open, setOpen] = useState(false);

    const handleTopUp = async () => {
        // Валидация на клиенте
        const value = parseFloat(amount);
        if (isNaN(value) || value <= 0) {
            setMessage({isError: true, message: "Введите корректную сумму"});
            return;
        }

        setLoading(true);
        setMessage(null);

        // Return URL - текущая страница
        const returnUrl = window.location.href;

        // Вызов API
        const url = await InitPayment(value, returnUrl, setMessage);

        if (url) {
            console.log("Redirecting to:", url);
            window.location.href = url;
        } else {
            // Если url нет, значит была ошибка, она уже в setMessage
            setLoading(false);
        }
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button
                    variant="outline"
                    size="sm"
                    className="ml-2 gap-2 border-[#DB935B] text-[#DB935B] hover:bg-[#DB935B] hover:text-black font-bold uppercase transition-all z-10 relative"
                >
                    <PlusCircle className="w-4 h-4" /> Пополнить
                </Button>
            </DialogTrigger>

            <DialogContent className="sm:max-w-[400px] bg-[#1e0e31] text-white border-gray-700">
                <DialogHeader>
                    <DialogTitle className="text-xl font-bold uppercase text-white flex items-center gap-2">
                        <CreditCard className="text-[#825e9c]"/>
                        Пополнение баланса
                    </DialogTitle>
                </DialogHeader>

                <div className="grid gap-6 py-4">
                    <div className="space-y-2">
                        <Label htmlFor="topup-amount" className="text-gray-300">Сумма пополнения (₽)</Label>
                        <div className="relative">
                            <Input
                                id="topup-amount"
                                type="number"
                                value={amount}
                                min={10}
                                onChange={(e) => {
                                    setAmount(e.target.value);
                                    setMessage(null);
                                }}
                                className="bg-white text-black border-gray-300 text-lg font-bold pl-4 pr-12 h-12"
                                placeholder="1000"
                            />
                            <span className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-500 font-bold">₽</span>
                        </div>
                    </div>

                    <div className="bg-[#825e9c]/10 border border-[#825e9c]/30 rounded p-3 text-xs text-gray-300">
                        <p>Вы будете перенаправлены на защищенный шлюз ЮKassa.</p>
                        <p className="mt-1">Комиссия сервиса: <b>0%</b></p>
                    </div>
                </div>

                <MessageComponent message={message} className="text-sm mb-2" />

                <Button
                    onClick={handleTopUp}
                    disabled={loading}
                    className="w-full bg-[#825e9c] hover:bg-[#6a4c80] text-white font-bold h-12 text-lg uppercase shadow-lg"
                >
                    {loading ? <Spinner size={24} /> : `Оплатить ${amount} ₽`}
                </Button>
            </DialogContent>
        </Dialog>
    );
}