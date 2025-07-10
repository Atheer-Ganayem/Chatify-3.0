import { Card } from "@/components/ui/card";
import useFetch from "@/hooks/useFetch";
import { Loader2, Trash2 } from "lucide-react";
import Image from "next/image";
import React from "react";
import { toast } from "sonner";

interface Props {
  url: string;
  isLoading: boolean;
  reset: () => void;
}

const ImagePreviewer: React.FC<Props> = ({ url, isLoading, reset }) => {
  const { exec, isLoading: isDeleting } = useFetch({
    path: "/image",
    method: "DELETE",
    auth: true,
    defaultLoading: false,
  });

  async function onDeleteHandler() {
    try {
      const { ok, responseData, error } = await exec();
      if (error) throw error;
      if (ok) {
        toast.success(responseData.message);
        reset();
      } else {
        toast.error(responseData.message);
      }
    } catch (error) {
      console.log(error);
      toast.error("Something went wrong, couldn't delete image.");
    }
  }

  return (
    <Card className="px-4 py-2 rounded-b-none">
      <div className="relative w-fit">
        <Image
          alt="img"
          src={url}
          width={100}
          height={0}
          style={{ height: "auto" }}
          className={isLoading || isDeleting ? "opacity-20" : undefined}
        />
        {(isLoading || isDeleting) && (
          <Loader2 className="animate-spin absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2" />
        )}
        {!isLoading && !isDeleting && (
          <div
            onClick={async () => await onDeleteHandler()}
            className="absolute top-0.5 right-0.5 text-destructive p-2 text-sm rounded-md backdrop-brightness-50 hover:cursor-pointer hover:backdrop-brightness-[25%]"
          >
            <Trash2 size={16} />
          </div>
        )}
      </div>
    </Card>
  );
};

export default ImagePreviewer;
