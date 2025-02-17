"use client";
import { AspectRatio } from "@qttdev/ui/aspect-ratio";
import "@qttdev/ui/styles.css";
import dynamic from "next/dynamic";

const loader = () =>
  import("@qttdev/ui/theme-picker").then((mod) => ({
    default: mod.ThemePicker,
  }));

const ThemePicker = dynamic(loader, { ssr: false });

export default function Home() {
  return (
    <>
      <ThemePicker />
      <main className="container mx-auto p-10 flex flex-col gap-5 min-h-screen bg-background">
        <h3 className="text-primary text-3xl cursor-pointer">Tizizz，一个简单的私人网络工具。</h3>
        <div className="bg-background p-5 rounded-md gap-5">
          <h5 className="text-primary/50 text-2xl cursor-pointer mt-5">{"> "}使用方式</h5>
          <p className="mt-5">微信扫描下方二维码关注 “huxulm”，获取使用教程。</p>
          <div className="w-64 mt-5">
            <AspectRatio ratio={1/1} className="size-full object-cover">
              <img src="./qrcode_for_gh.jpg"></img>
            </AspectRatio>
          </div>
        </div>
        <footer className="absolute bottom-10">
          <p className="font-song text-xs">Created by Huxulm © 2019-2025 Tizizz-top</p>
        </footer>
      </main>
    </>
  );
}
