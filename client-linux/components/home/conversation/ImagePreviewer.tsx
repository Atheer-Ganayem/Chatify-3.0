import { Card } from "@/components/ui/card";
import { Loader2 } from "lucide-react";
import Image from "next/image";
import React from "react";

interface Props {
  url: string;
  isLoading: boolean;
}

const ImagePreviewer: React.FC<Props> = ({ url, isLoading }) => {
  return (
    <Card className="px-4 py-2 rounded-b-none">
      <div className="relative w-fit">
        <Image
          alt="img"
          src={url}
          width={100}
          height={0}
          style={{ height: "auto" }}
          className={isLoading ? "opacity-20" : undefined}
        />
        {isLoading && (
          <Loader2 className="animate-spin absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2" />
        )}
      </div>
    </Card>
  );
};

export default ImagePreviewer;
