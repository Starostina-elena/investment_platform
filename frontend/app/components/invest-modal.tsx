'use client'
import {useState} from "react";
import {Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger} from "@/app/components/ui/dialog";
import {Button} from "@/app/components/ui/button";
import {Input} from "@/app/components/ui/input";
import {Label} from "@/app/components/ui/label";
import {MakeTransfer} from "@/api/transactions"; // Используем новый метод
import MessageComponent from "@/app/components/message";
import {Message} from "@/api/api";
import Spinner from "@/app/components/spinner";
import {useUserStore} from "@/context/user-store";
import Link from "next/link";
import {Wallet} from "lucide-react";
import TopUpModal from "@/app/components/top-up-modal";

export default function InvestModal({projectId, projectName, onSuccess}: {projectId: number, projectName: string, onSuccess?: () => void}) {
    const [amount, setAmount] = useState<string>("1000");
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState<Message | null>(null);
    const [open, setOpen] = useState(false);

    // Получаем текущего юзера и его баланс
    const { user } = useUserStore();
    const currentBalance = user?.balance || 0;
    const investAmount = parseFloat(amount) || 0;
    const isEnough = currentBalance >= investAmount;

    const handleInvest = async () => {
        if (!user) {
            setMessage({isError: true, message: "Необходимо войти в аккаунт"});
            return;
        }

        setLoading(true);
        setMessage(null);

        // Выполняем внутренний перевод
        const success = await MakeTransfer('project', projectId, investAmount, setMessage);

        setLoading(false);
        if (success) {
            setTimeout(() => {
                setOpen(false);
                if (onSuccess) onSuccess();
            }, 1500);
        }
    };

    if (!user) {
        // Если не залогинен - просто ссылка на вход
        return (
            <Link href="/login" className="text-white bg-[#00D0FF] font-bold text-center w-[40%] rounded-[40px] p-3 hover:brightness-110">
                Поддержать
            </Link>
        )
    }

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <a href="#" onClick={(e) => e.preventDefault()}
                   className="text-white bg-[#00D0FF] font-bold text-center w-[40%] rounded-[40px] p-3 hover:brightness-110">
                    Поддержать
                </a>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px] bg-[#656662] text-white border-gray-700">
                <DialogHeader>
                    <DialogTitle>Поддержать проект "{projectName}"</DialogTitle>
                </DialogHeader>

                <div className="grid gap-6 py-4">
                    {/* Баланс пользователя */}
                    <div className="flex items-center justify-between p-3 bg-white/5 rounded-lg border border-white/10">
                        <div className="flex items-center gap-2">
                            <Wallet className="text-[#DB935B]" />
                            <span>Мой кошелек:</span>
                        </div>
                        <span className="text-xl font-bold">{currentBalance} ₽</span>
                    </div>

                    <div className="grid gap-2">
                        <Label htmlFor="amount">Сумма инвестиции (₽)</Label>
                        <Input
                            id="amount"
                            type="number"
                            value={amount}
                            min={10}
                            onChange={(e) => {
                                setAmount(e.target.value);
                                setMessage(null);
                            }}
                            className="bg-[#301EBD] border-none text-white text-lg"
                        />
                    </div>

                    {/* Логика отображения кнопок */}
                    {!isEnough ? (
                        <div className="text-center space-y-3">
                            <p className="text-red-400 text-sm">
                                Недостаточно средств на кошельке. <br/>
                                Необходимо пополнить на: <b>{investAmount - currentBalance} ₽</b>
                            </p>
                            {/* Вставляем компонент пополнения прямо сюда, но стилизуем под широкую кнопку */}
                            <div className="bg-white/5 p-4 rounded-lg">
                                <TopUpModal />
                                {/* TopUpModal рендерит кнопку-триггер. Пользователь нажмет на нее, откроется второе окно пополнения. */}
                            </div>
                        </div>
                    ) : (
                        <div className="bg-[#B7FF00]/10 p-3 rounded text-[#B7FF00] text-sm text-center">
                            Средства будут списаны с вашего внутреннего кошелька.
                        </div>
                    )}
                </div>

                <MessageComponent message={message} />

                <Button
                    onClick={handleInvest}
                    disabled={loading || !isEnough}
                    className="bg-[#00D0FF] text-black hover:bg-[#00b0d6] font-bold h-12 w-full"
                >
                    {loading ? <Spinner size={24} /> : `Инвестировать ${amount} ₽`}
                </Button>
            </DialogContent>
        </Dialog>
    );
}