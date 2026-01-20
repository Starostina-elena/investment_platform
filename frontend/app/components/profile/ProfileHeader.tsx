"use client";

import React from 'react';
import { Camera } from 'lucide-react';
import { Button } from '@/app/components/ui/button';
import styles from './profileHeader.module.css';

interface ProfileHeaderProps {
  coverImage?: string;
  onCoverUpload: (file: File) => void;
}

export function ProfileHeader({ coverImage, onCoverUpload }: ProfileHeaderProps) {
  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      onCoverUpload(file);
    }
  };

  return (
    <div className={styles.coverContainer}>
      <div className={styles.coverPattern}></div>
      <div className={styles.coverOverlay}>
        <div className={styles.buttonContainer}>
          <input
            type="file"
            accept="image/*"
            onChange={handleFileChange}
            className={styles.hiddenInput}
            id="cover-upload"
          />
          <label htmlFor="cover-upload">
            <Button variant="outline" className={styles.uploadButton} asChild>
              <span>
                <Camera className="w-4 h-4 mr-2" />
                Добавить фото обложки
              </span>
            </Button>
          </label>
        </div>
      </div>
    </div>
  );
}