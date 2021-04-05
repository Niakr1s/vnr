import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { SentenceService } from '../services/sentence.service';

@Component({
  selector: 'app-progress',
  templateUrl: './progress.component.html',
  styleUrls: ['./progress.component.css'],
})
export class ProgressComponent implements OnInit, OnDestroy {
  total!: number;
  current!: number;

  subs: Subscription[] = [];

  constructor(private sentenceService: SentenceService) {}

  ngOnInit(): void {
    this.subs.push(
      this.sentenceService.totalSentences$.subscribe({
        next: (total) => {
          this.total = total;
        },
      })
    );
    this.subs.push(
      this.sentenceService.currentIndex$.subscribe({
        next: (current) => {
          this.current = current;
        },
      })
    );
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
  }

  textForDisplay(): string {
    const toString = (num: number) => {
      return num.toString().padStart(2, '0');
    };

    const current = this.current + 1;
    const total = this.total;

    return `${toString(current)}/${toString(total)}`;
  }
}
