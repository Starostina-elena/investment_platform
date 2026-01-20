'use client'
import { ReactNode } from 'react';

export interface SimpleModalProps {
    isOpen: boolean;
    onClose: () => void;
    title: string;
    children: ReactNode;
}

export default function SimpleModal({ isOpen, onClose, title, children }: SimpleModalProps) {
    if (!isOpen) return null;

    return (
        <>
            {/* Overlay */}
            <div
                className="fixed inset-0 bg-black/80 z-40"
                onClick={onClose}
            />
            {/* Modal */}
            <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
                <div className="bg-[#656662] text-white rounded-lg shadow-xl max-w-sm w-full max-h-[90vh] overflow-y-auto">
                    {/* Header */}
                    <div className="flex justify-between items-center p-4 border-b border-gray-600">
                        <h2 className="text-lg font-bold">{title}</h2>
                        <button
                            onClick={onClose}
                            className="text-gray-400 hover:text-white text-2xl leading-none"
                        >
                            ×
                        </button>
                    </div>
                    {/* Content */}
                    <div className="p-4">
                        {children}
                    </div>
                </div>
            </div>
        </>
    );
}
