.PHONY: all clean wiskcache ccache

subdirs = libX

# Generated Makefile will have many lib subdirs
#subdirs = ...

all wiskcache ccache: $(subdirs)

clean:
	@for i in $(subdirs); do 	(cd $$i; $(MAKE) clean); done


.PHONY: $(subdirs)
$(subdirs):
	PREFIX=$(PREFIX) $(MAKE) -C $@ $(MAKECMDGOALS)

depend:
	@echo Skip depend

demo_clean:
	@rm -rf /nobackup/$(USER)/wisk-nfscache/*

demo:  demo_clean
	@echo
	-@$(MAKE) clean >/dev/null 2>&1
	@echo
	@echo Regular build
	@/usr/bin/time  --format "%U;%S;%E;%PCPU" $(MAKE) all
	@echo ============================================
	@echo
	-@$(MAKE) clean >/dev/null 2>&1
	@echo
	@echo "Wiskcache build #1"
	@/usr/bin/time --format "%U;%S;%E;%PCPU"  $(MAKE) wiskcache
	@echo ============================================
	@echo
	-@$(MAKE) clean >/dev/null 2>&1
	@echo
	@echo "Wiskcache build #2"
	@/usr/bin/time --format "%U;%S;%E;%PCPU"  $(MAKE) wiskcache
	@echo ============================================
	@echo
	-@$(MAKE) clean >/dev/null 2>&1
	@echo
	@echo "Ccache build #1"
	@/usr/bin/time --format "%U;%S;%E;%PCPU"  $(MAKE) ccache > /dev/null
	@echo ============================================
	@echo
	-@$(MAKE) clean >/dev/null 2>&1
	@echo
	@echo "Ccache build #2"
	@/usr/bin/time --format "%U;%S;%E;%PCPU"  $(MAKE) ccache > /dev/null
	@echo ============================================
# DO NOT DELETE
