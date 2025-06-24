"use client";

import { ImageIcon, Plus } from "lucide-react";
import { FormControl, FormField, FormItem, FormMessage } from "../ui/form";
import { Input } from "../ui/input";
import { Control, FieldValues, Path } from "react-hook-form";
import { useRef, useState } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";

interface Props<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
  current?: string;
}

export default function AvatarInput<T extends FieldValues>({
  control,
  name,
  current,
}: Props<T>) {
  const [previewUrl, setPreviewUrl] = useState<string | undefined>(current);
  const inputRef = useRef<HTMLInputElement>(null);

  return (
    <>
      <Avatar
        className="h-24 w-24 mx-auto cursor-pointer"
        onClick={() => inputRef.current?.click()}
      >
        <AvatarImage src={previewUrl} alt="@shadcn" />
        <AvatarFallback>
          <Plus />
        </AvatarFallback>
      </Avatar>

      <FormField
        control={control}
        name={name}
        render={({ field }) => (
          <FormItem>
            <FormControl>
              <div className="relative" hidden>
                <ImageIcon className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  ref={inputRef}
                  type="file"
                  className="pl-10"
                  accept="image/*"
                  onChange={(e) => {
                    const file = e.target.files?.[0];
                    if (file) {
                      setPreviewUrl(URL.createObjectURL(file));
                    } else {
                      setPreviewUrl(undefined);
                    }

                    field.onChange(e.target.files);
                  }}
                />
              </div>
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
    </>
  );
}
